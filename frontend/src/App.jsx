import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import Navigation from './components/Navigation/Navigation'
import Register from './pages/Register'
import Login from './pages/Login'
import Authenticate from './pages/authenicate/Authenticate'

function App() {
  return (
    <Router>
      <Navigation />
      <Routes>
        <Route path="/" element={<Home />} />
        {/* <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} /> */}
        <Route path="/authenticate" element={<Authenticate />} />
      </Routes>
    </Router>
  )
}

export default App
