import { createSlice } from '@reduxjs/toolkit'


const initialState = {
    isAuth: false,
    user: null,
    otp: {
        phone: "",
        hash: "",
        expiresAt: "",
    }
}

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setAuth: (state, action) => {
            state.isAuth = true;
            state.user = action.payload;
        },

        setOtp: (state, action) => {
            state.otp = action.payload;
        },
        logoutAuth: (state) => {
            state.isAuth = false;
            state.user = null;
        }
    },

})

export const { setAuth, setOtp, logoutAuth } = authSlice.actions

export default authSlice.reducer