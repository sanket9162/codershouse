import { io } from "socket.io-client";

export const socketInit = () => {
    const options = {
        "force new connection": true,
        reconnectionAttempt: "Infinity",
        timeout: 10000,
        transports: ["websocket"],
    };

    const apiUrl = window.__ENV__?.API_URL || import.meta.env.VITE_API_URL;
    return io(apiUrl, options);
};
