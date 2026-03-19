import { FormEvent, useEffect, useState } from "react";
import { Tenant, getCurrentTenant, updateCurrentTenant } from "../../lib/api";

type TenantPanelProps = {
  token: string;
};

export function TenantPanel({ token }: TenantPanelProps) {
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getCurrentTenant(token)
      .then((nextTenant) => {
        if (cancelled) {
          return;
        }
        setTenant(nextTenant);
        setName(nextTenant.name);
      })
      .catch((err) => {
        if (cancelled) {
          return;
        }
        setMessage(err instanceof Error ? err.message : "failed to load tenant");
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [token]);

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!tenant) {
      return;
    }
    setMessage("");
    setLoading(true);
    try {
      const updated = await updateCurrentTenant(token, { name });
      setTenant(updated);
      setName(updated.name);
      setMessage("Tenant updated.");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "failed to update tenant");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section style={{ border: "1px solid #ccc", borderRadius: 8, padding: "1rem" }}>
      <h2 style={{ marginTop: 0 }}>Tenant Profile</h2>
      {loading && !tenant ? <p>Loading tenant...</p> : null}
      {tenant ? (
        <form onSubmit={onSubmit} style={{ display: "grid", gap: "0.5rem", maxWidth: 420 }}>
          <label>
            Tenant Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Tenant Slug
            <input value={tenant.slug} disabled />
          </label>
          <button disabled={loading || name.trim() === ""} type="submit">
            {loading ? "Saving..." : "Update tenant"}
          </button>
        </form>
      ) : null}
      {message ? <p style={{ marginBottom: 0 }}>{message}</p> : null}
    </section>
  );
}
