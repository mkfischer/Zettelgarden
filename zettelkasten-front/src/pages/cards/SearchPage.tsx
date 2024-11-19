import React, { useState, useEffect, ChangeEvent, KeyboardEvent } from "react";
import { fetchCards, semanticSearchCards } from "../../api/cards";
import { fetchUserTags } from "../../api/tags";
import { CardChunk, Card, PartialCard } from "../../models/Card";
import { Tag } from "../../models/Tags";
import { sortCards } from "../../utils/cards";
import { Button } from "../../components/Button";
import { CardList } from "../../components/cards/CardList";
import { CardChunkList } from "../../components/cards/CardChunkList";
import { SearchTagMenu } from "../../components/tags/SearchTagMenu";
import { usePartialCardContext } from "../../contexts/CardContext";

interface SearchPageProps {
  searchTerm: string;
  setSearchTerm: (searchTerm: string) => void;
  cards: PartialCard[];
  setCards: (cards: PartialCard[]) => void;
}

export function SearchPage({
  searchTerm,
  setSearchTerm,
  cards,
  setCards,
}: SearchPageProps) {
  const [sortBy, setSortBy] = useState("sortNewOld");
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(20);
  const { partialCards } = usePartialCardContext();
  const [useClassicSearch, setUseClassicSearch] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [chunks, setChunks] = useState<CardChunk[]>([]);

  const [tags, setTags] = useState<Tag[]>([]);

  function handleSearchUpdate(e: ChangeEvent<HTMLInputElement>) {
    setSearchTerm(e.target.value);
  }

  async function handleSearch(inputTerm = "") {
    setIsLoading(true);
    setCards([]);
    let term = inputTerm == "" ? searchTerm : inputTerm;

    try {
      if (useClassicSearch) {
        const data = await fetchCards(term);
        if (data === null) {
          setCards([]);
        } else {
          setCards(data);
        }
      } else {
        if (term === "") {
          setIsLoading(false);
          return;
        }
        const data = await semanticSearchCards(term);
        if (data === null) {
          setChunks([]);
        } else {
          setChunks(data);
        }
      }
    } catch (error) {
      console.error("Search error:", error);
      // Handle error appropriately
    } finally {
      setIsLoading(false);
    }
  }

  function handleSortChange(e: ChangeEvent<HTMLSelectElement>) {
    setSortBy(e.target.value);
  }

  function getSortedAndPagedCards() {
    const sortedCards = sortCards(cards, sortBy);
    const indexOfLastItem = currentPage * itemsPerPage;
    const indexOfFirstItem = indexOfLastItem - itemsPerPage;
    return sortedCards.slice(indexOfFirstItem, indexOfLastItem);
  }

  async function fetchTags() {
    fetchUserTags().then((data) => {
      if (data !== null) {
        console.log("tags");
        console.log(data);
        setTags(data);
      }
    });
  }

  function handleTagClick(tagName: string) {
    setSearchTerm("#" + tagName);
    handleSearch(tagName);
  }

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const term = params.get("term");
    if (term) {
      setSearchTerm(term);
      handleSearch(term);
    } else {
      fetchTags();
      handleSearch();
    }
    document.title = "Zettelgarden - Search";
  }, []);

  const currentItems = getSortedAndPagedCards();

  const handleCheckboxChange = (event) => {
    setUseClassicSearch(event.target.checked);
  };

  return (
    <div>
      <div>
        <div className="bg-slate-200 p-2 border-slate-400 border">
          <input
            style={{ display: "block", width: "100%", marginBottom: "10px" }} // Updated style here
            type="text"
            id="title"
            value={searchTerm}
            placeholder="Search"
            onChange={handleSearchUpdate}
            onKeyPress={(event: KeyboardEvent<HTMLInputElement>) => {
              if (event.key === "Enter") {
                handleSearch();
              }
            }}
          />

          <div className="flex">
            <Button onClick={() => handleSearch()} children={"Search"} />
            <select value={sortBy} onChange={handleSortChange}>
              <option value="sortNewOld">Newest</option>
              <option value="sortOldNew">Oldest</option>
              <option value="sortBigSmall">A to Z</option>
              <option value="sortSmallBig">Z to A</option>
            </select>
            <SearchTagMenu
              tags={tags.filter((tag) => tag.card_count > 0)}
              handleTagClick={handleTagClick}
            />
            <label>
              <input
                type="checkbox"
                checked={useClassicSearch}
                onChange={handleCheckboxChange}
              />
              Use Classic Search
            </label>
          </div>
        </div>
        {isLoading ? (
          <div className="flex justify-center w-full py-20">Loading</div>
        ) : (
          <div>
            {currentItems.length > 0 || chunks.length > 0 ? (
              <div>
                {useClassicSearch ? (
                  <CardList
                    cards={currentItems}
                    sort={false}
                    showAddButton={false}
                  />
                ) : (
                  <CardChunkList
                    cards={chunks}
                    sort={false}
                    showAddButton={false}
                  />
                )}
                <div>
                  <Button
                    onClick={() => setCurrentPage(currentPage - 1)}
                    disabled={currentPage === 1}
                    children={"Previous"}
                  />
                  <span>
                    {" "}
                    Page {currentPage} of{" "}
                    {Math.ceil(
                      (cards.length > 0 ? cards.length : partialCards.length) /
                        itemsPerPage,
                    )}{" "}
                  </span>
                  <Button
                    onClick={() => setCurrentPage(currentPage + 1)}
                    disabled={
                      currentPage ===
                      Math.ceil(
                        (cards.length > 0
                          ? cards.length
                          : partialCards.length) / itemsPerPage,
                      )
                    }
                    children={"Next"}
                  />
                </div>
              </div>
            ) : (
              <div className="flex justify-center w-full py-20">
                Search returned no results
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
