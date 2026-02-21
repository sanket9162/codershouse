import { useDispatch } from "react-redux";
import { useEffect, useState } from "react";
import axios from "axios";
import { setAuth } from '../store/authSlice';

export function useLoadingWithRefresh() {

    const [loading, setLoading] = useState(true);
    const dispatch = useDispatch();

    useEffect(() => {
        (async () => {
            try {
                const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/refresh-token`, {
                    withCredentials: true,
                });

                if (data && data.auth) {
                    dispatch(setAuth(data.user))

                }
                dispatch(setUser(data));

            } catch (error) {
                console.log(error)
            } finally {
                setLoading(false)
            }
        })()
    }, []);

    return { loading }
}