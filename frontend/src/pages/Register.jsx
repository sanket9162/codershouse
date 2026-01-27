import React from 'react'
import { useState } from 'react'
import StepPhoneEmail from './Steps/StepPhoneEmail'
import StepOtp from './Steps/StepOtp'
import StepName from './Steps/StepName'
import StepAvatar from './Steps/StepAvatar'
import StepUsername from './Steps/StepUsername'

const steps = {
    1: StepPhoneEmail,
    2: StepOtp,
    3: StepName,
    4: StepAvatar,
    5: StepUsername,
}

const Register = () => {
    const [step, setStep] = useState(3)
    const Step = steps[step]
    return (
        <div>
            <Step />
        </div>
    )
}



export default Register