const API_URL = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || '/api';
const API_PORT = process.env.API_PORT || '8080';

export function fetchApi(path: string, options?: RequestInit, version = 'v1'): Promise<any> {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  const url = `http://localhost:${API_PORT}${API_URL}/${version}${cleanPath}`;

  console.log(`Fetching API: ${url}`);

  return fetch(url, options).then((res) => {
    if (!res.ok) {
      throw new Error(`API request failed with status ${res.status}`);
    }
    return res.json();
  });
}


