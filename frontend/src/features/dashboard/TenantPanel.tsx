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
    <section className="card">
      <h2 className="panel-title">Tenant Profile</h2>
      {loading && !tenant ? <p>Loading tenant...</p> : null}
      {tenant ? (
        <form className="form-grid narrow" onSubmit={onSubmit}>
          <label>
            Tenant Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Tenant Slug
            <input value={tenant.slug} disabled />
          </label>
          <button className="btn-primary" disabled={loading || name.trim() === ""} type="submit">
            {loading ? "Saving..." : "Update tenant"}
          </button>
        </form>
      ) : null}
      {message ? (
        <p className={`message ${message.toLowerCase().includes("failed") ? "message-error" : "message-success"}`}>{message}</p>
      ) : null}
    </section>
  );
}
