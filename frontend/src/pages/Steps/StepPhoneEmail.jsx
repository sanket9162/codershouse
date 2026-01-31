import React, { useState } from 'react'
import Phone from './PhoneEmail/Phone'
import Email from './PhoneEmail/Email'

const phoneEmailMap = {
    phone: Phone,
    email: Email,
}

const StepPhoneEmail = ({onNext}) => {
  const [type, setType] = useState('phone');
  const Component = phoneEmailMap[type];

  return (
    <>

    <div className='flex items-center justify-center mt-4'> 
      <div>
      <div className='flex'>
        <button onClick={() => setType('phone')}>Phone</button>
        <button onClick={() => setType('email')}>Email</button>
      </div>
</div>
      <Component onNext={onNext}/>
    </div>
    </>
)
}

export default StepPhoneEmail