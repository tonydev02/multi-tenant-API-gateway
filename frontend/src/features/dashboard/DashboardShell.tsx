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
    <main style={{ fontFamily: "sans-serif", margin: "2rem auto", maxWidth: 1080 }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1.25rem" }}>
        <div>
          <h1 style={{ margin: 0 }}>Gateway Admin Dashboard</h1>
          <p style={{ marginTop: "0.5rem" }}>
            Signed in as {session.claims.email} (tenant #{session.claims.tenant_id})
          </p>
        </div>
        <button onClick={onLogout}>Logout</button>
      </header>

      <section style={{ display: "grid", gap: "1rem" }}>
        <TenantPanel token={session.token} />
        <ApiKeysPanel token={session.token} onChanged={onRefresh} />
        <TrafficPanel token={session.token} refreshTick={refreshTick} />
      </section>
    </main>
  );
}
