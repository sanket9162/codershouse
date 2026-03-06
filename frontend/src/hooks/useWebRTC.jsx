import { useState } from "react";
export const useWebRTC = () => {
    const [clients, setClients] = useState([
        {
            id: 1,
            name: 'Sanket',
        },
        {
            id: 2,
            name: 'gondhali',
        }
    ])
    return { clients }
}