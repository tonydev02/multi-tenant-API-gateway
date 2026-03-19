export type LoginResponse = { token: string };

export type Claims = {
  sub: number;
  tenant_id: number;
  email: string;
  iss: string;
  iat: number;
  exp: number;
};

export type Tenant = {
  id: number;
  name: string;
  slug: string;
  created_at: string;
  updated_at: string;
};

export type APIKey = {
  id: number;
  tenant_id: number;
  name: string;
  prefix: string;
  created_at: string;
  revoked_at?: string;
};

export type APIKeyCreateResult = {
  id: number;
  tenant_id: number;
  name: string;
  prefix: string;
  api_key: string;
  created_at: string;
};

export type TrafficSummary = {
  tenant_id: number;
  total_requests: number;
  rate_limited_requests: number;
  status_2xx: number;
  status_4xx: number;
  status_5xx: number;
  avg_latency_ms: number;
};

const apiBase = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${apiBase}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    },
    ...init
  });

  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(body.error ?? `request failed: ${res.status}`);
  }

  return (await res.json()) as T;
}

export async function registerTenant(input: {
  tenant_name: string;
  tenant_slug: string;
  email: string;
  password: string;
}): Promise<void> {
  await request("/api/admin/tenants/register", {
    method: "POST",
    body: JSON.stringify(input)
  });
}

export async function login(input: { email: string; password: string }): Promise<LoginResponse> {
  return request<LoginResponse>("/api/admin/login", {
    method: "POST",
    body: JSON.stringify(input)
  });
}

export async function getMe(token: string): Promise<Claims> {
  return request<Claims>("/api/admin/me", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`
    }
  });
}

function authHeaders(token: string): HeadersInit {
  return {
    Authorization: `Bearer ${token}`
  };
}

export async function getCurrentTenant(token: string): Promise<Tenant> {
  return request<Tenant>("/api/admin/tenants/current", {
    method: "GET",
    headers: authHeaders(token)
  });
}

export async function updateCurrentTenant(token: string, input: { name: string }): Promise<Tenant> {
  return request<Tenant>("/api/admin/tenants/current", {
    method: "PATCH",
    headers: authHeaders(token),
    body: JSON.stringify(input)
  });
}

export async function createAPIKey(token: string, input: { name: string }): Promise<APIKeyCreateResult> {
  return request<APIKeyCreateResult>("/api/admin/api-keys", {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify(input)
  });
}

export async function listAPIKeys(token: string): Promise<APIKey[]> {
  return request<APIKey[]>("/api/admin/api-keys", {
    method: "GET",
    headers: authHeaders(token)
  });
}

export async function revokeAPIKey(token: string, id: number): Promise<{ status: string }> {
  return request<{ status: string }>(`/api/admin/api-keys/${id}/revoke`, {
    method: "POST",
    headers: authHeaders(token)
  });
}

export async function getTrafficSummary(token: string): Promise<TrafficSummary> {
  return request<TrafficSummary>("/api/admin/traffic/summary", {
    method: "GET",
    headers: authHeaders(token)
  });
}
