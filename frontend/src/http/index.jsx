import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
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

export default api;
