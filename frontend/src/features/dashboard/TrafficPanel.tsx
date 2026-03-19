import { useEffect, useState } from "react";
import { TrafficSummary, getTrafficSummary } from "../../lib/api";

type TrafficPanelProps = {
  token: string;
  refreshTick: number;
};

export function TrafficPanel({ token, refreshTick }: TrafficPanelProps) {
  const [summary, setSummary] = useState<TrafficSummary | null>(null);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setMessage("");
    getTrafficSummary(token)
      .then((nextSummary) => {
        if (cancelled) {
          return;
        }
        setSummary(nextSummary);
      })
      .catch((err) => {
        if (cancelled) {
          return;
        }
        setSummary(null);
        setMessage(err instanceof Error ? err.message : "failed to load traffic summary");
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [token, refreshTick]);

  return (
    <section style={{ border: "1px solid #ccc", borderRadius: 8, padding: "1rem" }}>
      <h2 style={{ marginTop: 0 }}>Traffic Summary</h2>
      {loading ? <p>Loading traffic summary...</p> : null}
      {!loading && summary ? (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, minmax(0, 1fr))", gap: "0.5rem" }}>
          <Metric label="Total requests" value={summary.total_requests} />
          <Metric label="Rate-limited" value={summary.rate_limited_requests} />
          <Metric label="Avg latency (ms)" value={summary.avg_latency_ms} />
          <Metric label="2xx" value={summary.status_2xx} />
          <Metric label="4xx" value={summary.status_4xx} />
          <Metric label="5xx" value={summary.status_5xx} />
        </div>
      ) : null}
      {message ? <p style={{ marginBottom: 0 }}>{message}</p> : null}
    </section>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div style={{ border: "1px solid #ddd", borderRadius: 6, padding: "0.5rem" }}>
      <div style={{ fontSize: "0.8rem", color: "#444" }}>{label}</div>
      <div style={{ fontSize: "1.25rem", fontWeight: 700 }}>{value}</div>
    </div>
  );
}
