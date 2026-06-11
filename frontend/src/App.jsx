import { useState } from 'react'
import './App.css'
import {Route,Routes} from "react-router-dom"
import Home from "./pages/Home"
import Report from './pages/Report'
import User from './pages/User'
function App() {
  return (
    <Routes>
        <Route path="/" element={<Home/>}></Route>
         <Route path="/user" element={<User/>}></Route>
        <Route path="/report" element={<Report/>}></Route>
    </Routes>
  )
}

export default App
