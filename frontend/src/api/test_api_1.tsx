const url = "http://localhost:8000/test/8";

import type { NextApiRequest, NextApiResponse } from "next";

export default async function testAPI1(
	req: NextApiRequest,
	res: NextApiResponse
) {
	const data = await fetch(url, {
		method: "GET",
		headers: {
			"Content-Type": "application/json",
		},
	});
    console.log(data)
    const word = await data.json()
    console.log(word)
    res.status(200).json({
        message: word,
    })
}
