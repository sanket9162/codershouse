import axios from "axios";

const api = axios.create({
    baseURL: window.__ENV__?.API_URL || import.meta.env.VITE_API_URL,
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
    },
});

export const sendOTP = (phone) => {
    return api.post("/send-otp", { phone });
}

export const verifyOTP = (phone, otp, hash, expiresAt) => {
    return api.post("/verify-otp", { phone, otp, hash, expiresAt });
}

export const activateUser = (name, avatar) => {
    return api.post("/activate", { name, avatar });
}

export const logout = () => {
    return api.get("/logout")
}

export const createRoom = (topic, roomType) => {
    return api.post("/rooms", { topic, roomType });
}

export const getAllRooms = () => {
    return api.get("/rooms");
}

export const getRoomById = (roomId) => {
    return api.get(`/rooms/${roomId}`);
}

export const searchRoomsAPI = (query) => api.get(`/rooms?search=${query}`);

// interceptors
api.interceptors.response.use(
    (config) => {
        return config;
    },
    async (error) => {
        const originalRequest = error.config;
        if (error.response.status === 401 &&
            originalRequest &&
            !originalRequest._isRetry
        ) {
            originalRequest._isRetry = true;
            try {
                const apiUrl = window.__ENV__?.API_URL || import.meta.env.VITE_API_URL;
                await axios.get(`${apiUrl}/refresh`, {
                    withCredentials: true,
                });
                return api.request(originalRequest);
            } catch (error) {
                return console.log("Refresh token failed");
            }
        }
        throw error;
    }
)

export default api;
