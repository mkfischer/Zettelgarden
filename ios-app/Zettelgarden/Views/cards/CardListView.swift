import SwiftUI
import ZettelgardenShared

struct CardListView: View {
    @State private var errorMessage: String?
    @State private var isInitialized: Bool = false
    @ObservedObject var cardViewModel: CardViewModel
    @ObservedObject var navigationViewModel: NavigationViewModel
    @ObservedObject var viewModel: PartialCardViewModel

    var body: some View {
        VStack {
            FilterFieldView(filterText: $viewModel.filterText, placeholder: "Filter")
            if viewModel.isLoading {
                ProgressView("Loading")
            }
            else if let _ = viewModel.cards {

                List {
                    ForEach(viewModel.filteredCards) { card in
                        CardListItem(
                            card: card,
                            navigationViewModel: navigationViewModel
                        )
                    }
                }
                .refreshable {
                    viewModel.loadCards()
                }
            }
        }
    }
}
