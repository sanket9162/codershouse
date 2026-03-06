import { useStateWithCallback } from "./useStateWithCallback";

const users = [
    {
        id: 1,
        name: 'Sanket',
    },
    {
        id: 2,
        name: 'gondhali',
    }
]

export const useWebRTC = () => {
    const [clients, setClients] = useStateWithCallback(users)
    return { clients }
}