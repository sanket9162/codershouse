import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { useSelector } from 'react-redux'

const SemiProtectedRoute = ({ children, ...rest }) => {
    const isAuth = useSelector((state) => state.auth.isAuth);
    const isActivated = useSelector((state) => state.auth.user.activated);

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
