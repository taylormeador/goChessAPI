import requests
import db
import time

LICHESS_API_URL = "https://explorer.lichess.ovh/masters"
GO_CHESS_API_URL = "https://go-chess-api.herokuapp.com/isLegal"
START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
DEPTH = 10
BREADTH = 2


class InsertQueryValues:
    def __init__(self, fen, name=None):
        self.fen = fen
        self.name = name
        self.moves = None
        self.move_1 = None
        self.move_2 = None
        self.move_3 = None
        self.move_4 = None
        self.move_5 = None
        self.move_6 = None
        self.move_7 = None
        self.move_8 = None
        self.move_9 = None
        self.move_10 = None

    def populate_moves(self):
        for i in range(len(self.moves)):
            if i == 0:
                self.move_1 = self.moves[i]
            elif i == 1:
                self.move_2 = self.moves[i]
            elif i == 2:
                self.move_3 = self.moves[i]
            elif i == 3:
                self.move_4 = self.moves[i]
            elif i == 4:
                self.move_5 = self.moves[i]
            elif i == 5:
                self.move_6 = self.moves[i]
            elif i == 6:
                self.move_7 = self.moves[i]
            elif i == 7:
                self.move_8 = self.moves[i]
            elif i == 8:
                self.move_9 = self.moves[i]
            elif i == 9:
                self.move_10 = self.moves[i]


# query the lichess api endpoint and return the json response
def get_request(fen, tries=0):
    url = LICHESS_API_URL + "?fen=" + fen
    try:
        r = requests.get(url)
    except Exception:
        tries += 1
        if tries <= 10:
            time.sleep(5)
            get_request(fen, tries)
        else:
            r = "Error getting fen from Lichess"
    return r.json()


value_objects = []

# parse the json response and extract the top n moves
def get_top_moves(fen, moves):
    json = get_request(fen)
    print(f"fen: {fen}")
    try:
        name = json['opening']['name']
        values = InsertQueryValues(fen, name)
    except TypeError:
        values = InsertQueryValues(fen)

    move_list = []
    for i in range(moves):
        # check if there is a move, if not continue
        try:
            uci_move = json['moves'][i]['uci']
        except IndexError:
            continue
        new_fen = make_move(fen, uci_move)['FEN']
        move_list.append(new_fen)

    values.moves = move_list
    value_objects.append(values)
    return move_list


# calls the Go Chess API to return the fen after making the move
def make_move(fen, uci_format_move):
    url = GO_CHESS_API_URL + "?FEN=" + fen + " moves " + uci_format_move
    r = requests.get(url)
    return r.json()


# recursively deepen move tree
def driver(fen, depth):
    if depth == 0:
        return

    moves = get_top_moves(fen, BREADTH)
    for new_fen in moves:
        driver(new_fen, depth - 1)


def insert_in_db(values):
    print("*********************************************")
    cur, conn = db.connect()
    # check if the entry already exists
    query = f"""SELECT * FROM opening_book WHERE fen = (%s);"""
    print(f"Selecting from db: {values.fen}")
    cur.execute(query, (values.fen,))
    result = cur.fetchone()

    # check if the entry already exists
    if result:
        print(f"FEN already exists in db, checking breadth....")
        num_of_moves = 10 - result[2:].count(None)
        if num_of_moves < len(values.moves):
            # TODO update the move_n col
            print("Additional moves available")
            return
        print("No new moves available - nothing to do")
        return
    else:  # add the entry if it doesn't exist
        # query
        query = f"""INSERT INTO opening_book (fen, name, move_1, move_2, move_3, move_4, move_5, 
        move_6, move_7, move_8, move_9, move_10) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);"""

        print(f"Inserting to db: {values.fen}, {values.name}, {values.moves}")
        cur.execute(query, (values.fen, values.name, values.move_1, values.move_2, values.move_3, values.move_4,
                            values.move_5, values.move_6, values.move_7, values.move_8, values.move_9, values.move_10))

        # Make the changes to the database persistent
        conn.commit()


# main func call
driver(START_FEN, DEPTH)
for value_object in value_objects:
    value_object.populate_moves()
    insert_in_db(value_object)


def print_results():
    cur, conn = db.connect()
    # query
    query = f"""SELECT * FROM opening_book;"""
    cur.execute(query)
    print(cur.fetchall(), '\n')


# print_results()  # debug
