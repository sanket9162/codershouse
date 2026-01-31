import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'

const GuestRoute = ({ children }) => {
    const isAuth = true;
    return isAuth ? <Navigate to="/rooms" /> : (children ? children : <Outlet />);
}

export default GuestRoute
