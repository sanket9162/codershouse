import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { useSelector } from 'react-redux'

const ProtectedRoute = ({ children }) => {
    const isAuth = useSelector((state) => state.auth.isAuth);
    const isActivated = useSelector((state) => state.auth.user.activated);

    if (!isAuth) {
        return <Navigate to="/" />
    }

    if (!isActivated) {
        return <Navigate to="/activate" />
    }

    return children ? children : <Outlet />;
}

export default ProtectedRoute
