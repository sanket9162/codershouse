import React from 'react'
import { Link, useNavigate } from 'react-router-dom'
import Card from '../components/Card/Card'
import Button from '../components/Button/Button'

export default function Home() {
  const navigate = useNavigate()
  const startRegister = () => {
    navigate('/authenticate')
  }
  
  return (
    <div className="flex items-center justify-center">
      <Card title="Welcome to Coderhouse" icon="logo">
       <p className="text-lg leading-relaxed">
            We&apos;re working hard to get Coderhouse ready for everyone! While 
            we wrap up the finishing touches, we&apos;re adding people gradually
            to make sure nothing breaks.
          </p>
          <div>
            <Button onClick={startRegister} text="Let's Go"/>
          </div>
          <div className="mt-4">
            <span className="text-[#0077ff]">Have an invite text?</span>
          </div>
      </Card>
    </div>
  )
}
