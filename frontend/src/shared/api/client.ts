export async function apiClient(
  url: string,
  options: RequestInit = {},
): Promise<Response> {
  const sendRequest = async (accessToken: string) => {
    fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
    });
  };

  let response = await sendRequest(getAccessToken());

  if (response.status === 401) {
    const tokens = await refreshAccessToken();
    response = await sendRequest(getAccessToken());
  }

  return response;
}

function getAccessToken() {
  return localStorage.getItem("accessToken");
}

function getRefreshToken() {
  return localStorage.getItem("refreshToken");
}

async function refreshAccessToken() {
  return;
}
