import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
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
                await axios.get(`${import.meta.env.VITE_API_URL}/refresh`, {
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
