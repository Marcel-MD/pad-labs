import { check, sleep } from 'k6'
import http from 'k6/http'

import jsonpath from 'https://jslib.k6.io/jsonpath/1.0.2/index.js'

import { randomIntBetween, randomString, randomItem } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: '1m30s', target: 50 },
        { duration: '30s', target: 5 },
    ],
};

export default function main() {
    let users = [
        {
            email: 'marcel@gmail.com',
            name: 'Marcel',
            password: 'password',
        },
        {
            email: 'masha@gmail.com',
            name: 'Masha',
            password: 'password',
        },
        {
            email: 'bostan@gmail.com',
            name: 'Bostan',
            password: 'password',
        },
    ]

    let response

    let randomUser = randomItem(users)
    let userBody = JSON.stringify(randomUser)

    response = http.post(
        'http://localhost:3010/users/login',
        userBody,
        {
            headers: {
                'content-type': 'application/json',
                accept: 'application/json',
                'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
                'sec-ch-ua-mobile': '?0',
                'sec-ch-ua-platform': '"Linux"',
            },
        }
    )

    if (response.status > 299) {
        response = http.post(
            'http://localhost:3010/users/register',
            userBody,
            {
                headers: {
                    'content-type': 'application/json',
                    accept: 'application/json',
                    'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
                    'sec-ch-ua-mobile': '?0',
                    'sec-ch-ua-platform': '"Linux"',
                },
            }
        )
    }

    let token = jsonpath.query(response.json(), '$.token')[0]

    sleep(1)

    response = http.get('http://localhost:3010/users/current', {
        headers: {
            authorization: `Bearer ${token}`,
            accept: 'application/json',
            'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
            'sec-ch-ua-mobile': '?0',
            'sec-ch-ua-platform': '"Linux"',
        },
    })

    check(response, {
        'is status 200': (r) => r.status === 200 || r.status === 201,
        'is name correct': (r) => jsonpath.query(r.json(), '$.name')[0] === randomUser.name,
    })

    sleep(1)

    let randomProduct = {
        name: randomString(10),
        price: randomIntBetween(100, 5000),
        stock: randomIntBetween(1, 100)
    }

    let productBody = JSON.stringify(randomProduct)

    response = http.post(
        'http://localhost:3010/product',
        productBody,
        {
            headers: {
                authorization: `Bearer ${token}`,
                'content-type': 'application/json',
                accept: 'application/json',
                'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
                'sec-ch-ua-mobile': '?0',
                'sec-ch-ua-platform': '"Linux"',
            },
        }
    )

    check(response, {
        'is status 200': (r) => r.status === 200 || r.status === 201,
        'id returned': (r) => jsonpath.query(r.json(), '$.id')[0] !== undefined,
    })

    let productId = jsonpath.query(response.json(), '$.id')[0]

    sleep(2)

    response = http.get('http://localhost:3010/product', {
        headers: {
            accept: 'application/json',
            'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
            'sec-ch-ua-mobile': '?0',
            'sec-ch-ua-platform': '"Linux"',
        },
    })

    check(response, {
        'is status 200': (r) => r.status === 200
    })

    sleep(2)

    let randomOrder = {
        shippingAddress: randomString(10),
        quantity: randomIntBetween(1, randomProduct.stock)
    }

    let orderBody = JSON.stringify(randomOrder)

    response = http.post(
        `http://localhost:3010/order/${productId}`,
        orderBody,
        {
            headers: {
                authorization: `Bearer ${token}`,
                'content-type': 'application/json',
                accept: 'application/json',
                'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
                'sec-ch-ua-mobile': '?0',
                'sec-ch-ua-platform': '"Linux"',
            },
        }
    )

    check(response, {
        'is status 200': (r) => r.status === 200 || r.status === 201,
        'id returned': (r) => jsonpath.query(r.json(), '$.id')[0] !== undefined,
    })

    sleep(2)

    response = http.get('http://localhost:3010/order', {
        headers: {
            accept: 'application/json',
            'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
            'sec-ch-ua-mobile': '?0',
            'sec-ch-ua-platform': '"Linux"',
        },
    })

    check(response, {
        'is status 200': (r) => r.status === 200
    })
}