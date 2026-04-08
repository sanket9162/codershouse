import React, { useState } from 'react'
import Button from '../../../components/Button/Button'
import Card from '../../../components/Card/Card'
import TextInput from '../../../components/TextInput/TextInput'

const Email = ({ onNext }) => {
  const [email, setEmail] = useState('')
  return (
    <Card title="Enter Email Address" icon="email-emoji">
      <TextInput value={email} onChange={(e) => setEmail(e.target.value)} />
      <div className='flex flex-col items-center justify-center'>
        <div className='mt-2'>
          <Button text="Next" onClick={onNext} />
        </div>
        <p className='w-4/5 text-sm mt-4 text-gray-400'>By entering your number, you&apos;re agreeing to our Terms of Service and Privacy Policy</p>
      </div>
    </Card>
  )
}

export default Email