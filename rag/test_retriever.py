from retriever import search


def test_fresh_fish_query_prefers_market_candidates():
    results = search("川崎で2000円以内、新鮮な魚")
    assert results
    assert results[0]["name"] in {"さか本", "鮨あらい"}


def test_budget_filter():
    results = search("川崎 魚", budget_yen_max=2000)
    assert results
    assert all(r["metadata"]["budget_yen"] <= 2000 for r in results)


def test_area_filter():
    results = search("魚", area="川崎")
    assert results
    assert all(r["metadata"]["area"] == "川崎" for r in results)
