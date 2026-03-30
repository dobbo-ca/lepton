import { A } from "@solidjs/router";
import { Component } from "solid-js";

interface NavItem {
  href: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { href: "/", label: "Dashboard", icon: "⊞" },
  { href: "/agents", label: "Agents", icon: "⚡" },
  { href: "/issues", label: "Issues", icon: "◎" },
  { href: "/runs", label: "Runs", icon: "▶" },
  { href: "/routines", label: "Routines", icon: "⏱" },
  { href: "/settings", label: "Settings", icon: "⚙" },
];

const Sidebar: Component = () => {
  return (
    <nav
      class="flex flex-col w-56 shrink-0 bg-slate-900 border-r border-slate-800 h-full"
      aria-label="Main navigation"
    >
      <div class="px-5 py-4 border-b border-slate-800">
        <span class="text-slate-100 font-semibold text-sm tracking-wide">lepton</span>
      </div>
      <ul class="flex flex-col gap-0.5 p-2 mt-1" role="list">
        {navItems.map((item) => (
          <li>
            <A
              href={item.href}
              class="flex items-center gap-3 px-3 py-2 rounded text-sm text-slate-400 hover:text-slate-100 hover:bg-slate-800 transition-colors"
              activeClass="bg-slate-800 text-slate-100"
              end={item.href === "/"}
              aria-label={item.label}
            >
              <span class="w-4 text-center" aria-hidden="true">
                {item.icon}
              </span>
              {item.label}
            </A>
          </li>
        ))}
      </ul>
    </nav>
  );
};

export default Sidebar;
