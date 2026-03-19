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
    <section className="card">
      <h2 className="panel-title">Traffic Summary</h2>
      {loading ? <p>Loading traffic summary...</p> : null}
      {!loading && summary ? (
        <div className="metrics-grid">
          <Metric label="Total requests" value={summary.total_requests} />
          <Metric label="Rate-limited" value={summary.rate_limited_requests} />
          <Metric label="Avg latency (ms)" value={summary.avg_latency_ms} />
          <Metric label="2xx" value={summary.status_2xx} />
          <Metric label="4xx" value={summary.status_4xx} />
          <Metric label="5xx" value={summary.status_5xx} />
        </div>
      ) : null}
      {message ? <p className="message message-error">{message}</p> : null}
    </section>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="metric-card">
      <div className="metric-label">{label}</div>
      <div className="metric-value">{value}</div>
    </div>
  );
}
