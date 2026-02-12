import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { useSelector } from 'react-redux'

const GuestRoute = ({ children }) => {
    const isAuth = useSelector((state) => state.auth.isAuth);
    return isAuth ? <Navigate to="/rooms" /> : (children ? children : <Outlet />);
}

export default GuestRoute
