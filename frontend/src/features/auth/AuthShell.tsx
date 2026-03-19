import { FormEvent, useState } from "react";
import { Claims, getMe, login, registerTenant } from "../../lib/api";

type Mode = "login" | "register";

type AuthShellProps = {
  onAuthenticated: (session: { token: string; claims: Claims }) => void;
};

export function AuthShell({ onAuthenticated }: AuthShellProps) {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("admin@acme.local");
  const [password, setPassword] = useState("changeme123");
  const [tenantName, setTenantName] = useState("Acme");
  const [tenantSlug, setTenantSlug] = useState("acme");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setMessage("");
    setLoading(true);

    try {
      if (mode === "register") {
        await registerTenant({
          tenant_name: tenantName,
          tenant_slug: tenantSlug,
          email,
          password
        });
        setMessage("Tenant registered. You can now login.");
        setMode("login");
      } else {
        const auth = await login({ email, password });
        const me = await getMe(auth.token);
        onAuthenticated({ token: auth.token, claims: me });
      }
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="auth-shell">
      <section className="card hero-card">
        <h1 className="hero-title">Gateway Admin Dashboard</h1>
        <p className="hero-subtitle">
          Production-style multi-tenant API gateway SaaS built with Go for concurrent, network-heavy backend
          workloads.
        </p>
        <p className="hero-details">
          This admin UI demonstrates secure tenant onboarding, JWT-based admin authentication, API key lifecycle
          management, tenant-aware request controls, and live traffic visibility for REST gateway operations.
        </p>
        <p className="hero-link">
          GitHub repository:{" "}
          <a href="https://github.com/tonydev02/multi-tenant-API-gateway" rel="noreferrer" target="_blank">
            github.com/tonydev02/multi-tenant-API-gateway
          </a>
        </p>
      </section>

      <section className="card">
        <div className="mode-switch">
          <button
            className={mode === "login" ? "btn-primary" : "btn-secondary"}
            disabled={loading || mode === "login"}
            onClick={() => setMode("login")}
            type="button"
          >
            Login
          </button>
          <button
            className={mode === "register" ? "btn-primary" : "btn-secondary"}
            disabled={loading || mode === "register"}
            onClick={() => setMode("register")}
            type="button"
          >
            Register Tenant
          </button>
        </div>

        <form className="form-grid spaced-top" onSubmit={onSubmit}>
          {mode === "register" && (
            <>
              <label>
                Tenant Name
                <input value={tenantName} onChange={(e) => setTenantName(e.target.value)} />
              </label>
              <label>
                Tenant Slug
                <input value={tenantSlug} onChange={(e) => setTenantSlug(e.target.value)} />
              </label>
            </>
          )}
          <label>
            Email
            <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" />
          </label>
          <label>
            Password
            <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" />
          </label>
          <button className="btn-primary" disabled={loading} type="submit">
            {loading ? "Working..." : mode === "register" ? "Create tenant" : "Login"}
          </button>
        </form>

        {message && (
          <p
            className={`message ${
              message.toLowerCase().includes("error") || message.toLowerCase().includes("failed")
                ? "message-error"
                : "message-success"
            }`}
          >
            {message}
          </p>
        )}
      </section>
    </main>
  );
}
