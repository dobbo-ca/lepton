import { Component, JSX } from "solid-js";

interface TopBarProps {
  title: string;
  breadcrumbs?: { label: string; href?: string }[];
  actions?: JSX.Element;
}

const TopBar: Component<TopBarProps> = (props) => {
  return (
    <header class="flex items-center justify-between px-6 py-3 border-b border-slate-800 bg-slate-950 shrink-0">
      <div class="flex flex-col gap-0.5">
        {props.breadcrumbs && props.breadcrumbs.length > 0 && (
          <nav aria-label="Breadcrumb">
            <ol class="flex items-center gap-1 text-xs text-slate-500">
              {props.breadcrumbs.map((crumb, i) => (
                <>
                  {i > 0 && <span aria-hidden="true">/</span>}
                  <li>
                    {crumb.href ? (
                      <a href={crumb.href} class="hover:text-slate-300 transition-colors">
                        {crumb.label}
                      </a>
                    ) : (
                      <span class="text-slate-400">{crumb.label}</span>
                    )}
                  </li>
                </>
              ))}
            </ol>
          </nav>
        )}
        <h1 class="text-sm font-medium text-slate-100">{props.title}</h1>
      </div>
      {props.actions && <div class="flex items-center gap-2">{props.actions}</div>}
    </header>
  );
};

export default TopBar;
