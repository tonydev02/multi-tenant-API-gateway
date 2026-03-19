import { useState } from "react";
import { Claims } from "../../lib/api";
import { ApiKeysPanel } from "./ApiKeysPanel";
import { TenantPanel } from "./TenantPanel";
import { TrafficPanel } from "./TrafficPanel";

type DashboardShellProps = {
  session: { token: string; claims: Claims };
  onLogout: () => void;
};

export function DashboardShell({ session, onLogout }: DashboardShellProps) {
  const [refreshTick, setRefreshTick] = useState(0);

  function onRefresh() {
    setRefreshTick((value) => value + 1);
  }

  return (
    <main>
      <header className="card top-panel dashboard-header">
        <div>
          <h1 className="hero-title">Gateway Admin Dashboard</h1>
          <p className="dashboard-meta">
            Signed in as {session.claims.email} (tenant #{session.claims.tenant_id})
          </p>
          <p className="dashboard-link">
            Project source:{" "}
            <a href="https://github.com/tonydev02/multi-tenant-API-gateway" rel="noreferrer" target="_blank">
              GitHub
            </a>
          </p>
        </div>
        <button className="btn-secondary" onClick={onLogout}>
          Logout
        </button>
      </header>

      <section className="dashboard-grid">
        <section className="card app-overview">
          <h2 className="panel-title">What This Web App Can Do</h2>
          <ul className="feature-list">
            <li>Authenticate platform admins with JWT-backed sessions.</li>
            <li>Manage tenant profile details used by gateway services.</li>
            <li>Create and revoke tenant-scoped API keys for consumers.</li>
            <li>Surface tenant traffic metrics including rate-limited requests and latency.</li>
            <li>Provide an operational control plane for a multi-tenant REST API gateway.</li>
          </ul>
          <p className="feature-repo-link">
            Explore implementation details on{" "}
            <a href="https://github.com/tonydev02/multi-tenant-API-gateway" rel="noreferrer" target="_blank">
              GitHub
            </a>
            .
          </p>
        </section>
        <TenantPanel token={session.token} />
        <ApiKeysPanel token={session.token} onChanged={onRefresh} />
        <TrafficPanel token={session.token} refreshTick={refreshTick} />
      </section>
    </main>
  );
}
