import { FormEvent, useState } from "react";
import { Claims, getMe, login, registerTenant } from "../../lib/api";

type Mode = "login" | "register";

export function AuthShell() {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("admin@acme.local");
  const [password, setPassword] = useState("changeme123");
  const [tenantName, setTenantName] = useState("Acme");
  const [tenantSlug, setTenantSlug] = useState("acme");
  const [token, setToken] = useState("");
  const [claims, setClaims] = useState<Claims | null>(null);
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
        setToken(auth.token);
        const me = await getMe(auth.token);
        setClaims(me);
        setMessage("Login successful.");
      }
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={{ fontFamily: "sans-serif", margin: "2rem auto", maxWidth: 960 }}>
      <h1>Gateway Admin Dashboard</h1>
      <p>Phase 02 auth and tenancy shell.</p>

      <div style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem" }}>
        <button disabled={loading || mode === "login"} onClick={() => setMode("login")}>
          Login
        </button>
        <button disabled={loading || mode === "register"} onClick={() => setMode("register")}>
          Register Tenant
        </button>
      </div>

      <form onSubmit={onSubmit} style={{ display: "grid", gap: "0.75rem", maxWidth: 420 }}>
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
        <button disabled={loading} type="submit">
          {loading ? "Working..." : mode === "register" ? "Create tenant" : "Login"}
        </button>
      </form>

      {message && <p style={{ marginTop: "1rem" }}>{message}</p>}

      {claims && (
        <section style={{ marginTop: "1.5rem" }}>
          <h2>Session</h2>
          <ul>
            <li>Email: {claims.email}</li>
            <li>Tenant ID: {claims.tenant_id}</li>
            <li>Admin ID: {claims.sub}</li>
          </ul>
          <details>
            <summary>JWT</summary>
            <code>{token}</code>
          </details>
        </section>
      )}
    </main>
  );
}
