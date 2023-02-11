import React, { useCallback, useEffect, useState } from 'react';

const socket = new WebSocket('ws://127.0.0.1:3000/ws');

function App() {
    const [message, setMessage] = useState('');
    const [inputValue, setInputValue] = useState('');
    const [data, setData] = useState([]);
    const [showPopup, setShowPopup] = useState(false);
    const [isFetching, setIsFetching] = useState(false);

    useEffect(() => {
        socket.onopen = () => {
            console.log('Connected');
            setMessage('Connected');
        };

        socket.onmessage = (e) => {
            const response = JSON.parse(e.data);
            setData(response.data);
        };

        socket.onclose = () => {
            console.log('Disconnected');
            setMessage('Disconnected');
        };

        try {
            socket.send(
                JSON.stringify({
                    q: inputValue,
                })
            );
        } catch (err) {
            console.error(err);
        }
        return () => {};
    }, [inputValue]);

    const handleClick = useCallback(
        (e) => {
            e.preventDefault();

            try {
                socket.send(
                    JSON.stringify({
                        message: inputValue,
                    })
                );
            } catch (err) {
                console.error(err);
            }
        },
        [inputValue]
    );

    const handleChange = useCallback((e) => {
        setInputValue(e.target.value);
    }, []);

    const handleClose = useCallback((e) => {
        e.preventDefault();
        socket.close();
    }, []);

    const togglePopup = () => {
        setShowPopup(!showPopup);
    };

    const highlightKeyword = (text, keyword) => {
        if (keyword === '') return text;
        const regex = new RegExp(`(${keyword})`, 'gi');
        return text.replace(regex, '<span class="bg-yellow-500">$1</span>');
    };

    return (
        <main className='bg-black min-h-screen w-full p-12'>
            <div className='relative max-w-lg mx-auto'>
                <input
                    value={inputValue}
                    type='text'
                    className='w-full p-3 rounded-lg shadow-sm border border-gray-500 bg-gray-800 text-gray-100 placeholder-gray-500 focus:border-gray-500 focus:outline-none'
                    placeholder='Search...'
                    onChange={handleChange}
                    onFocus={togglePopup}
                    onBlur={togglePopup}
                />

                <div className='absolute right-0 top-0 mt-3 mr-3'>
                    <button className='text-gray-600 hover:text-gray-700 focus:outline-none'>
                        <svg
                            className='w-5 h-5'
                            fill='none'
                            strokeLinecap='round'
                            strokeLinejoin='round'
                            strokeWidth={2}
                            viewBox='0 0 24 24'
                            stroke='currentColor'
                        >
                            <path d='M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z' />
                        </svg>
                    </button>
                </div>

                {showPopup && (
                    <div className='bg-white mt-2 rounded-lg shadow-lg absolute z-10 w-full max-h-80 overflow-y-auto'>
                        <ul className='list-reset'>
                            {data.length != 0 &&
                                data.map((item, index) => (
                                    <li key={index}>
                                        <a
                                            href='#'
                                            className='block px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white'
                                        >
                                            <div
                                                dangerouslySetInnerHTML={{
                                                    __html: highlightKeyword(
                                                        item.name,
                                                        inputValue
                                                    ),
                                                }}
                                            />
                                        </a>
                                    </li>
                                ))}
                        </ul>
                    </div>
                )}
            </div>
        </main>
    );
}

export default App;
