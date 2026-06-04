import { useState } from 'react'
import './App.css'
import {Route,Routes} from "react-router-dom"
import Home from "./pages/Home"
import Report from './pages/Report'
function App() {
  return (
    <Routes>
        <Route path="/" element={<Home/>}></Route>

        <Route path="/report" element={<Report/>}></Route>
    </Routes>
  )
}

export default App
