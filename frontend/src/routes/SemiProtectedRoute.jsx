import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'

const SemiProtectedRoute = ({ children, ...rest }) => {
    const isAuth = false; // validation
    const isActivated =    false; // validation

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
