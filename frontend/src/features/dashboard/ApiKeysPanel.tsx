import { FormEvent, useEffect, useState } from "react";
import { APIKey, APIKeyCreateResult, createAPIKey, listAPIKeys, revokeAPIKey } from "../../lib/api";

type ApiKeysPanelProps = {
  token: string;
  onChanged: () => void;
};

export function ApiKeysPanel({ token, onChanged }: ApiKeysPanelProps) {
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [newKeyName, setNewKeyName] = useState("default");
  const [created, setCreated] = useState<APIKeyCreateResult | null>(null);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);

  async function refreshKeys() {
    const nextKeys = await listAPIKeys(token);
    setKeys(nextKeys);
  }

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    refreshKeys()
      .catch((err) => {
        if (cancelled) {
          return;
        }
        setMessage(err instanceof Error ? err.message : "failed to load api keys");
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

  async function onCreate(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setMessage("");
    setCreated(null);
    setLoading(true);
    try {
      const result = await createAPIKey(token, { name: newKeyName });
      setCreated(result);
      await refreshKeys();
      onChanged();
      setMessage("API key created. Copy the secret now.");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "failed to create api key");
    } finally {
      setLoading(false);
    }
  }

  async function onRevoke(id: number) {
    setMessage("");
    setLoading(true);
    try {
      await revokeAPIKey(token, id);
      await refreshKeys();
      onChanged();
      setMessage("API key revoked.");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "failed to revoke api key");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="card">
      <h2 className="panel-title">API Keys</h2>
      <form className="inline-form" onSubmit={onCreate}>
        <input
          aria-label="API key name"
          value={newKeyName}
          onChange={(e) => setNewKeyName(e.target.value)}
          placeholder="Key name"
        />
        <button className="btn-primary" disabled={loading || newKeyName.trim() === ""} type="submit">
          {loading ? "Working..." : "Create key"}
        </button>
      </form>

      {created ? (
        <div className="info-box">
          <strong>New key (shown once):</strong>
          <pre>{created.api_key}</pre>
        </div>
      ) : null}

      {loading && keys.length === 0 ? <p>Loading keys...</p> : null}
      {keys.length > 0 ? (
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Prefix</th>
              <th>Created</th>
              <th>Revoked</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {keys.map((key) => (
              <tr key={key.id}>
                <td>{key.name}</td>
                <td>
                  <code>{key.prefix}</code>
                </td>
                <td>{new Date(key.created_at).toLocaleString()}</td>
                <td>{key.revoked_at ? new Date(key.revoked_at).toLocaleString() : "active"}</td>
                <td>
                  <button
                    className={key.revoked_at ? "btn-secondary" : "btn-danger"}
                    disabled={loading || !!key.revoked_at}
                    onClick={() => onRevoke(key.id)}
                    type="button"
                  >
                    Revoke
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <p>No API keys yet.</p>
      )}

      {message ? (
        <p
          className={`message ${
            message.toLowerCase().includes("failed") || message.toLowerCase().includes("error")
              ? "message-error"
              : "message-success"
          }`}
        >
          {message}
        </p>
      ) : null}
    </section>
  );
}
