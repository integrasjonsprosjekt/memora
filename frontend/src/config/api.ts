const getApiUrl = (): string => {
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
};

export const API_URL = getApiUrl();

// Helper function to construct API endpoints
export const getApiEndpoint = (path: string): string => {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  return `${API_URL}${cleanPath}`;
};

// TODO: Temporary workaround until we have a proper authentication system
const getUserId = (): string => {
  let user_id: string | undefined;

  if (typeof window !== 'undefined') {
    user_id = process.env.NEXT_PUBLIC_USER_ID;
  } else {
    user_id = process.env.USER_ID || process.env.NEXT_PUBLIC_USER_ID;
  }

  if (!user_id) {
    throw new Error('USER_ID is not defined.');
  }

  return user_id;
};

export const USER_ID = '1nhjfvpIl1ZTpxivPopOQLRbvrB3';
