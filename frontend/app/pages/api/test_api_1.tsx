import type { NextApiRequest, NextApiResponse } from 'next'

interface Test1 {
    x: string
}

const url = 'http://localhost:8000/1'

export default async function TestAPI1() {
    const res = await fetch(url, {
        method: 'GET',
        
    });
    
    const body = await res.json()
    console.log(body)
}
