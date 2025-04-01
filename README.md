# test_task_bkh_ekom_ethereum_converter
Test task for golang developer

## HOW TO USE 
```
go run main.go <Ethereum Address>
```

## RESULT
```
Balance in Ethereum:  0.000782810643493349
Balance in USD:       1.50
```

Задача 2. Расчет эквивалента USD на заданном адресе сети Ethereum
Нужен простой сервис, который сможет пересчитать кол-во ETH на заданном адресе в USD

Внешние АПИ использовать нельзя (только on-chain). Рекомендуется использовать контракты chainlink (например https://etherscan.io/address/0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419)

 

Для приложения не нужен веб интерфейс. Оно получает ethereum адрес в параметрах командной строки, выводит кол-во ETH на балансе и эквивалент этой суммы в USD.

 

Дополнительно: учесть WETH , если таковые есть на балансе

 

Макс. уровень сложности: вывести все токены и их эквивалент в USD, которые поддерживает chainlink. Только если они есть на балансе адреса (https://docs.chain.link/data-feeds/price-feeds/addresses?network=ethereum&page=3&search=USD)