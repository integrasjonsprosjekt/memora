import { User } from 'firebase/auth';

function getApiUrl(): string {
  let url: string | undefined;

  // For client-side components, use NEXT_PUBLIC_ prefix
  if (typeof window !== 'undefined') {
    url = process.env.NEXT_PUBLIC_API_URL;
  } else {
    url = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL;
  }

  if (!url) {
    console.warn('API_URL is not defined, using default');
    return '/api';
  }

  return url;
}

export async function fetchApi<T>(path: string, options?: RequestInit & { user?: User }, version = 'v1'): Promise<T> {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  const apiUrl = getApiUrl();
  const url = `${apiUrl}/${version}${cleanPath}`;

  // Extract user from options and remove it from the RequestInit object
  const { user, ...fetchOptions } = options || {};

  // Add Authorization header if user is provided
  if (user) {
    const idToken = await user.getIdToken();
    fetchOptions.headers = {
      ...fetchOptions.headers,
      Authorization: `Bearer ${idToken}`,
    };
  }

  return fetch(url, fetchOptions).then(async (res) => {
    if (!res.ok) {
      throw new Error(`API request failed with status ${res.status}`);
    }

    // Handle empty responses (common for DELETE requests with 204 No Content)
    const contentLength = res.headers.get('Content-Length');
    if (contentLength === '0' || res.status === 204) {
      return {} as T;
    }

    // Check if there's actual content to parse
    const text = await res.text();
    if (!text) {
      return {} as T;
    }

    return JSON.parse(text);
  });
}
