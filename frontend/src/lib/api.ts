export type LoginResponse = { token: string };

export type Claims = {
  sub: number;
  tenant_id: number;
  email: string;
  iss: string;
  iat: number;
  exp: number;
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
