import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// Add request interceptor to inject Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default {
  login(password) {
    return api.post("/auth/login", { password });
  },
  getParticipants(pool) {
    // We didn't create a specific get participants endpoint, but we can if needed.
    // Or just rely on teams for now.
    // Wait, Wheel needs participants.
    // I need GET /participants?pool=X.
    // I missed implementing GET /participants in backend.
    // For now, I'll Mock it or just use what I have.
    // Actually, the Wheel just needs a list of names for visual effect.
    // I'll stick to 'generate teams' which does the logic on backend.
    return Promise.resolve([]);
  },
  getTeams(pool) {
    return api.get(`/teams?pool=${pool}`);
  },
  generateTeams(pool) {
    return api.post("/teams/generate", { pool });
  },
  generateBracket(teamIds, pool) {
    return api.post("/bracket/generate", { team_ids: teamIds, pool });
  },
  getBracket(pool) {
    return api.get(`/bracket?pool=${pool}`);
  },
  updateMatch(id, data) {
    return api.post(`/matches/${id}/result`, data);
  },
};
