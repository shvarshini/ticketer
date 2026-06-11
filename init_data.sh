#!/bin/bash

# Exit on error
set -e

# API Base URL (defaults to localhost if no argument is provided)
BASE_URL=${1:-"http://localhost:8080/api"}
echo "Using BASE_URL: $BASE_URL"

# Helper function to extract ID from JSON using Python
extract_id() {
  python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))"
}

# Helper function to format JSON output
pretty_print() {
  python3 -m json.tool || cat
}

echo "==================================="
echo "1. Creating Theater..."
THEATER_RES=$(curl -s -X POST "$BASE_URL/admin/theaters" \
  -H "Content-Type: application/json" \
  -d '{
    "admin_id": "admin-123",
    "name": "PVR Cinemas",
    "location": "Downtown Mall"
  }')
echo "$THEATER_RES" | pretty_print
THEATER_ID=$(echo "$THEATER_RES" | extract_id)
echo "-> Extracted Theater ID: $THEATER_ID"

echo "==================================="
echo "2. Creating Screen..."
SCREEN_RES=$(curl -s -X POST "$BASE_URL/admin/theaters/$THEATER_ID/screens" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Screen 1 - IMAX"
  }')
echo "$SCREEN_RES" | pretty_print
SCREEN_ID=$(echo "$SCREEN_RES" | extract_id)
echo "-> Extracted Screen ID: $SCREEN_ID"

echo "==================================="
echo "3. Creating 6 Seats..."
for i in {1..6}; do
  SEAT_RES=$(curl -s -X POST "$BASE_URL/admin/screens/$SCREEN_ID/seats" \
    -H "Content-Type: application/json" \
    -d "{
      \"row\": \"A\",
      \"number\": $i,
      \"type\": \"NORMAL\"
    }")
  echo "Seat A-$i:"
  echo "$SEAT_RES" | pretty_print
done

echo "==================================="
echo "4. Creating Movie..."
MOVIE_RES=$(curl -s -X POST "$BASE_URL/admin/movies" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception",
    "description": "A thief who steals corporate secrets through the use of dream-sharing technology...",
    "duration": 148,
    "release_date": "2010-07-16T00:00:00Z",
    "genre": "Sci-Fi",
    "base_price": 15.50
  }')
echo "$MOVIE_RES" | pretty_print
MOVIE_ID=$(echo "$MOVIE_RES" | extract_id)
echo "-> Extracted Movie ID: $MOVIE_ID"

echo "==================================="
echo "5. Creating Show..."
SHOW_RES=$(curl -s -X POST "$BASE_URL/admin/shows" \
  -H "Content-Type: application/json" \
  -d "{
    \"movie_id\": \"$MOVIE_ID\",
    \"screen_id\": \"$SCREEN_ID\",
    \"start_time\": \"2026-10-15T18:00:00Z\",
    \"end_time\": \"2026-10-15T20:30:00Z\"
  }")
echo "$SHOW_RES" | pretty_print
SHOW_ID=$(echo "$SHOW_RES" | extract_id)
echo "-> Extracted Show ID: $SHOW_ID"

echo "==================================="
echo "6. Getting Show Seats..."
SHOW_SEATS_RES=$(curl -s -X GET "$BASE_URL/shows/$SHOW_ID/seats")
echo "$SHOW_SEATS_RES" | pretty_print

echo "==================================="
echo "Initialization complete!"
