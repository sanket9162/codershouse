import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'

const SemiProtectedRoute = ({ children, ...rest }) => {
    const isAuth = true; // validation
    const isActivated =    true; // validation

    if (!isAuth) {
        return <Navigate to="/" />
    }

    if (isAuth && !isActivated) {
        return children ? children : <Outlet />
    } else {
        return <Navigate to="/rooms" />
    }
}

export default SemiProtectedRoute
