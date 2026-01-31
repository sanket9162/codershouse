import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'

const ProtectedRoute = ({ children }) => {
    const isAuth = true; // validation
    const isActivated = true; // validation

    if (!isAuth) {
        return <Navigate to="/" />
    }

    if (!isActivated) {
        return <Navigate to="/activate" />
    }

    return children ? children : <Outlet />;
}

export default ProtectedRoute
