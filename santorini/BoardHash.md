# Implement a hash that effectively stores the game board

The board can be stored in a byte array very efficiently

Each tile gets 4 bits (level 0-2)
bit/team 1 | bit/team 2 | 2 bits/level (0-3)
Capp = team1 & team2
Empty = !team1 & !team2
