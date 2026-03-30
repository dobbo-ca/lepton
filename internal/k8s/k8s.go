// Package k8s provides the Kubernetes job runner and a local dev stub.
package k8s

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dobbo-ca/lepton/internal/domain"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// AgentRunner is the interface for submitting and managing agent executions.
type AgentRunner interface {
	Submit(ctx context.Context, agent domain.Agent, issue domain.Issue) (*domain.Run, error)
	Status(ctx context.Context, runID string) (domain.RunStatus, error)
	Cancel(ctx context.Context, runID string) error
	Logs(ctx context.Context, runID string) (io.Reader, error)
}

// ---- Kubernetes runner ----

// KubernetesRunner submits real Kubernetes Jobs via client-go.
type KubernetesRunner struct {
	client    kubernetes.Interface
	namespace string
}

// NewKubernetesRunner creates a KubernetesRunner using in-cluster config, falling back to kubeconfig.
func NewKubernetesRunner(namespace string) (*KubernetesRunner, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("k8s config: %w", err)
		}
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("k8s client: %w", err)
	}
	if namespace == "" {
		namespace = "default"
	}
	return &KubernetesRunner{client: client, namespace: namespace}, nil
}

func (k *KubernetesRunner) Submit(ctx context.Context, agent domain.Agent, issue domain.Issue) (*domain.Run, error) {
	now := time.Now()
	jobName := fmt.Sprintf("lepton-%s-%d", strings.ToLower(agent.Name), now.Unix())

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: k.namespace,
			Labels: map[string]string{
				"app":      "lepton",
				"agentId":  agent.ID,
				"issueId":  issue.ID,
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "agent",
							Image: agent.Image,
							Env: []corev1.EnvVar{
								{Name: "LEPTON_AGENT_ID", Value: agent.ID},
								{Name: "LEPTON_ISSUE_ID", Value: issue.ID},
							},
						},
					},
				},
			},
		},
	}

	if _, err := k.client.BatchV1().Jobs(k.namespace).Create(ctx, job, metav1.CreateOptions{}); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	run := &domain.Run{
		ID:        jobName,
		AgentID:   agent.ID,
		IssueID:   issue.ID,
		Status:    domain.RunStatusPending,
		StartedAt: &now,
	}
	return run, nil
}

func (k *KubernetesRunner) Status(ctx context.Context, runID string) (domain.RunStatus, error) {
	job, err := k.client.BatchV1().Jobs(k.namespace).Get(ctx, runID, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("get job: %w", err)
	}
	switch {
	case job.Status.Succeeded > 0:
		return domain.RunStatusSucceeded, nil
	case job.Status.Failed > 0:
		return domain.RunStatusFailed, nil
	case job.Status.Active > 0:
		return domain.RunStatusRunning, nil
	default:
		return domain.RunStatusPending, nil
	}
}

func (k *KubernetesRunner) Cancel(ctx context.Context, runID string) error {
	propagation := metav1.DeletePropagationForeground
	return k.client.BatchV1().Jobs(k.namespace).Delete(ctx, runID, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
}

func (k *KubernetesRunner) Logs(ctx context.Context, runID string) (io.Reader, error) {
	pods, err := k.client.CoreV1().Pods(k.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "job-name=" + runID,
	})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}
	if len(pods.Items) == 0 {
		return strings.NewReader(""), nil
	}
	req := k.client.CoreV1().Pods(k.namespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{})
	return req.Stream(ctx)
}

// ---- Local runner (dev stub) ----

// LocalRunner is an in-process stub for running agents without Kubernetes.
type LocalRunner struct {
	runs map[string]*domain.Run
}

// NewLocalRunner creates a LocalRunner.
func NewLocalRunner() *LocalRunner {
	return &LocalRunner{runs: make(map[string]*domain.Run)}
}

func (l *LocalRunner) Submit(_ context.Context, agent domain.Agent, issue domain.Issue) (*domain.Run, error) {
	now := time.Now()
	id := fmt.Sprintf("local-%s-%d", agent.ID, now.UnixNano())
	run := &domain.Run{
		ID:        id,
		AgentID:   agent.ID,
		IssueID:   issue.ID,
		Status:    domain.RunStatusRunning,
		StartedAt: &now,
	}
	l.runs[id] = run
	return run, nil
}

func (l *LocalRunner) Status(_ context.Context, runID string) (domain.RunStatus, error) {
	r, ok := l.runs[runID]
	if !ok {
		return "", errors.New("run not found: " + runID)
	}
	return r.Status, nil
}

func (l *LocalRunner) Cancel(_ context.Context, runID string) error {
	r, ok := l.runs[runID]
	if !ok {
		return errors.New("run not found: " + runID)
	}
	r.Status = domain.RunStatusCancelled
	return nil
}

func (l *LocalRunner) Logs(_ context.Context, runID string) (io.Reader, error) {
	if _, ok := l.runs[runID]; !ok {
		return nil, errors.New("run not found: " + runID)
	}
	return bytes.NewReader([]byte("(local runner — no logs)\n")), nil
}
