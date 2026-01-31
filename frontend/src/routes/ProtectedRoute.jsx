import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'

const ProtectedRoute = ({ children }) => {
    const isAuth = false; // validation
    const isActivated = false; // validation

    if (!isAuth) {
        return <Navigate to="/" />
    }

    if (!isActivated) {
        return <Navigate to="/activate" />
    }

    return children ? children : <Outlet />;
}

export default ProtectedRoute
