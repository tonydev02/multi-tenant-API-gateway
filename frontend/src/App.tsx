import { useMemo, useState } from "react";
import { AuthShell } from "./features/auth/AuthShell";
import { DashboardShell } from "./features/dashboard/DashboardShell";
import { Claims } from "./lib/api";

type Session = {
  token: string;
  claims: Claims;
};

const sessionKey = "gateway_admin_session";

function loadSession(): Session | null {
  const raw = window.sessionStorage.getItem(sessionKey);
  if (!raw) {
    return null;
  }
  try {
    return JSON.parse(raw) as Session;
  } catch {
    window.sessionStorage.removeItem(sessionKey);
    return null;
  }
}

export function App() {
  const [session, setSession] = useState<Session | null>(() => loadSession());

  const authenticated = useMemo(() => {
    if (!session) {
      return false;
    }
    return session.claims.exp * 1000 > Date.now();
  }, [session]);

  function onAuthenticated(nextSession: Session) {
    window.sessionStorage.setItem(sessionKey, JSON.stringify(nextSession));
    setSession(nextSession);
  }

  function onLogout() {
    window.sessionStorage.removeItem(sessionKey);
    setSession(null);
  }

  if (!authenticated || !session) {
    return <AuthShell onAuthenticated={onAuthenticated} />;
  }

  return <DashboardShell session={session} onLogout={onLogout} />;
}
