export const saveToken = (token: string) => {
  localStorage.setItem("piggy_token", token);
  document.cookie = `piggy_token=${token}; path=/; max-age=${60 * 60 * 24}; SameSite=Lax`;
};

export const getToken = (): string | null =>
  typeof window !== "undefined" ? localStorage.getItem("piggy_token") : null;

export const removeToken = () => {
  localStorage.removeItem("piggy_token");
  document.cookie = "piggy_token=; path=/; max-age=0";
};

export const authHeaders = (): Record<string, string> => {
  const token = getToken();
  return {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};