const API_URL = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || '/api';

export function fetchApi(path: string, options?: RequestInit, version = 'v1'): Promise<any> {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  const url = `${API_URL}/${version}${cleanPath}`;

  return fetch(url, options).then((res) => {
    if (!res.ok) {
      throw new Error(`API request failed with status ${res.status}`);
    }
    return res.json();
  });
}


